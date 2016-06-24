## finder_machine

Automatiza la transferencia de comprobantes entre diferentes orígenes de datos, como bases de datos, colas y sistema de archivos.

## Setup

    $ go get github.com/edmt/finder_machine

## Stored Procedures

FinderMachine_ReadXml:

```
CREATE PROCEDURE FinderMachine_ReadXml
    @startDate VARCHAR(20),
    @endDate varchar(20)
AS
BEGIN

SET NOCOUNT ON

select xml.uuid, xml.xml, xml.timestamp
from xml
inner join cfd on cfd.numTimbre = xml.uuid
left join POOL_ENVIOCFD_SAT_Z as pz on pz.comprobante_Id = cfd.idInternal
left join cfd_envio_sat_z as acuse on acuse.comprobante_Id = cfd.idInternal
where timestamp > @startDate and timestamp < @endDate and pz.comprobante_Id is null and acuse.comprobante_Id is null
order by timestamp desc

END


exec FinderMachine_ReadXML '2016-05-16', '2016-05-17'
```

FinderMachine_WritePool:

```
CREATE PROCEDURE FinderMachine_WritePool
    @uuid VARCHAR(40)
AS
BEGIN

SET NOCOUNT ON

declare @comprobanteId varchar(40)

set @comprobanteId = (select idinternal
                      from dbo.cfd
                      where numtimbre = @uuid)

insert into dbo.POOL_ENVIOCFD_SAT_Z(idInternal, comprobante_Id, fechaRegistro, status)
values(replace(newid(), '-', ''), @comprobanteId, getdate(), 0)

END

exec FinderMachine_WritePool 'fc0d9501-25f9-40fa-b4ba-73dfaf06dc6d'
```

FinderMachine_ReadXml_MissingCfd:

```
CREATE PROCEDURE FinderMachine_ReadXml_MissingCfd
    @startDate VARCHAR(20),
    @endDate varchar(20)
AS
BEGIN

SET NOCOUNT ON

select xml.uuid, xml.xml, xml.timestamp
from xml
left join cfd on cfd.numTimbre = xml.uuid
where timestamp > @startDate and timestamp < @endDate and cfd.idInternal is null
order by timestamp desc

END

exec FinderMachine_ReadXml_MissingCfd '2016-06-24', '2016-06-30'

```
