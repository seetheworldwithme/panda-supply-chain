// 离线提取原 HanLP portable-1.8.4 的公共 Apache-2.0 模型和对照结果，生产 Go 不依赖 JVM。
// javac -cp hanlp-portable-1.8.4.jar ExportHanLP.java
// java -cp .:hanlp-portable-1.8.4.jar ExportHanLP model.bin.gz golden.tsv
import java.io.*;
import java.lang.reflect.*;
import java.nio.charset.StandardCharsets;
import java.util.*;
import java.util.zip.GZIPOutputStream;
import com.hankcs.hanlp.HanLP;
import com.hankcs.hanlp.collection.trie.DoubleArrayTrie;
import com.hankcs.hanlp.corpus.dictionary.item.EnumItem;
import com.hankcs.hanlp.corpus.tag.*;
import com.hankcs.hanlp.dictionary.*;
import com.hankcs.hanlp.dictionary.nr.*;
import com.hankcs.hanlp.dictionary.other.*;
import com.hankcs.hanlp.seg.common.Term;
import com.hankcs.hanlp.utility.Predefine;

public class ExportHanLP {
  static Object field(Object target,String name)throws Exception{
    Class<?> c=target instanceof Class?(Class<?>)target:target.getClass();
    for(;c!=null;c=c.getSuperclass())try{Field f=c.getDeclaredField(name);f.setAccessible(true);return f.get(target instanceof Class?null:target);}catch(NoSuchFieldException e){}
    throw new NoSuchFieldException(name);
  }
  static void ints(DataOutputStream o,int[] a)throws Exception{o.writeInt(a.length);for(int v:a)o.writeInt(v);}
  static void str(DataOutputStream o,String s)throws Exception{byte[] b=s.getBytes(StandardCharsets.UTF_8);o.writeInt(b.length);o.write(b);}
  static void dat(DataOutputStream o,DoubleArrayTrie<?> t)throws Exception{
    int n=(Integer)field(t,"size");int[] b=(int[])field(t,"base"),c=(int[])field(t,"check");
    o.writeInt(n);for(int i=0;i<n;i++){o.writeInt(b[i]);o.writeInt(c[i]);}
  }
  static void attrs(DataOutputStream o,DoubleArrayTrie<?> t)throws Exception{
    Object[] values=(Object[])field(t,"v");o.writeInt(values.length);
    for(Object value:values){CoreDictionary.Attribute a=(CoreDictionary.Attribute)value;o.writeInt(a.totalFrequency);o.writeInt(a.nature.length);for(int i=0;i<a.nature.length;i++){o.writeInt(a.nature[i].ordinal());o.writeInt(a.frequency[i]);}}
    dat(o,t);
  }
  public static void main(String[] args)throws Exception{
    HanLP.segment("初始化");CustomDictionary.add("主焦");CustomDictionary.add("1/3焦");
    PersonDictionary.dictionary.get("张");TranslatedPersonDictionary.containsKey("·");CoreBiGramTableDictionary.getBiFrequency(0,0);
    try(DataOutputStream o=new DataOutputStream(new BufferedOutputStream(new GZIPOutputStream(new FileOutputStream(args[0]))))){
      o.write("HLP184".getBytes(StandardCharsets.US_ASCII));
      o.writeInt(Predefine.TOTAL_FREQUENCY);o.writeDouble(Predefine.myu);o.writeInt(Predefine.OOV_DEFAULT_FREQUENCY);o.writeBoolean(HanLP.Config.Normalization);
      Nature[] ns=Nature.values();o.writeInt(ns.length);for(Nature n:ns)str(o,n.toString());
      o.write(CharType.type);for(char c:CharTable.CONVERT)o.writeChar(c);
      attrs(o,CoreDictionary.trie);attrs(o,CustomDictionary.DEFAULT.dat);
      ints(o,(int[])field(CoreBiGramTableDictionary.class,"start"));ints(o,(int[])field(CoreBiGramTableDictionary.class,"pair"));
      DoubleArrayTrie<?> nr=(DoubleArrayTrie<?>)field(PersonDictionary.dictionary,"trie");Object[] items=(Object[])field(nr,"v");o.writeInt(items.length);
      for(Object v:items){Map<?,?> labels=((EnumItem<?>)v).labelMap;o.writeInt(labels.size());for(Map.Entry<?,?> e:labels.entrySet()){o.writeInt(((NR)e.getKey()).ordinal());o.writeInt((Integer)e.getValue());}}
      dat(o,nr);
      TransformMatrixDictionary matrix=PersonDictionary.transformMatrixDictionary;
      for(NR n:NR.values())o.writeInt(matrix.getTotalFrequency(n));
      for(double[] row:matrix.transititon_probability)for(double d:row)o.writeDouble(d);
      dat(o,(DoubleArrayTrie<?>)field(TranslatedPersonDictionary.class,"trie"));
    }
    String[] samples={"","主焦","1/3焦","主焦煤 低硫 0.6","山西主焦煤灰10硫0.8","山西省吕梁市孝义市主焦煤","临汾 焦煤 1000","山西大同动力煤5000大卡","张三 王小明 李四买煤","张三丰的煤矿","阿尔伯特·爱因斯坦","粗苯价格上涨","贫瘦煤低硫煤","1/3焦煤 0.80 10.00","全水分8 灰分10 挥发25","挥发分25%低硫0.6","CSR65 GR.I80 Y15","煤炭公司 价格1000","熊猫煤仓MC001","ABC 123","1,234.56 1.5e3 NaN Infinity","ＡＢＣ１２３ ０．８","精煤，原煤。焦炭！","主焦煤\n肥煤\t瘦煤","煤😀炭","水洗煤5000大卡硫0.6灰10","四川攀枝花供应焦煤","北京大学张华","乌鲁木齐煤矿","本周价格下降50元",
      "陕西榆林神木煤矿现货报价","内蒙古鄂尔多斯动力煤Q5500","晋城无烟煤块煤低硫","炼焦煤主焦1/3焦肥煤瘦煤","河北唐山钢厂采购焦炭","河南平顶山焦煤价格上涨","新疆哈密露天矿原煤","宁夏银川煤化工用煤","山东济宁气煤供应","云南曲靖焦煤质量指标","煤矿企业联系人张先生电话13800138000","水分≤8.5 灰分≤12.5 硫分≤0.6","发热量 5500 kcal/kg 挥发分 28%","灰份10.00，含硫量0.45，固定碳55","Q=5000大卡；S<0.8；V=25","G值75，Y值12，CSR60，CRI30","收到基低位发热量Qnet,ar=5500","空气干燥基灰分Ad 12.3%","全硫St,d不超过1.0%","粒度 0-50mm，块率大于80%","2026-11-19报价￥850.00/吨","价格从790元涨到810元，涨幅2.53%","联系人：李四（销售经理）","王建国、赵小明和刘海涛负责供应","张三丰先生来自北京市海淀区","欧阳娜娜采购俄罗斯煤炭","公司名称：山西焦煤集团有限责任公司","供货地址：山西省吕梁市孝义市兑镇镇","山西省长治市襄垣县煤矿资源","鄂尔多斯市东胜区罕台镇煤场","abcDEF ABCdef x86_64 HTTP/2 API v3","A1B2C3煤种XK-2026-001","12.5万吨，约3.6亿大卡","1/3焦煤、1/2中粘煤、气肥煤","无烟煤(三号) [优质] {现货}","煤炭/焦炭/兰炭/煤泥/煤矸石","矿区—站台—港口，铁路运输","联系电话(+86) 010-12345678","邮箱 sales@example.com 官网 https://coal.example.cn","价格区间：800~950元/吨（含税）","-10℃至+35℃，湿度≤80%","订单号 PO202611190001，合同 HT-001","规格型号MC-5500-Q/2026版","山西焦煤集团西山煤电（集团）有限责任公司","内蒙古伊泰煤炭股份有限公司","陕西煤业化工集团神木天元化学有限公司","中华人民共和国山西省吕梁市孝义市","第一名 第二名 第三名，第十次交易","今天天气晴朗，适合煤炭运输和装车","请问这个煤种的硫分和发热量是多少？","包头到天津港汽运运费每吨120元","神华准格尔能源有限责任公司黑岱沟露天煤矿","中煤平朔集团有限公司安太堡露天矿","低硫主焦煤（S≤0.6%）月度长协价","焦煤指数：JCI、CCI、CCTD、BSPI","6000±100大卡，收到基挥发份18.5-22.0","2025年12月31日之前有效，2026/01/01起执行","+86 139-0000-0000 / 400-123-4567","原煤洗选后精煤产率约72.5%，中煤15.0%","A股上市公司：601225.SH 陕西煤业","人民币(CNY) 1,000.00元；美元(USD) $150.00","🎉优质动力煤现货供应🔥欢迎咨询！"};
    try(PrintWriter w=new PrintWriter(new OutputStreamWriter(new FileOutputStream(args[1]),StandardCharsets.UTF_8))){
      for(String sample:samples){w.print(Base64.getEncoder().encodeToString(sample.getBytes(StandardCharsets.UTF_8)));for(Term t:HanLP.segment(sample)){w.print('\t');w.print(Base64.getEncoder().encodeToString(t.word.getBytes(StandardCharsets.UTF_8)));}w.println();}
    }
  }
}
